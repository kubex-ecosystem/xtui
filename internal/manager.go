// Package manager implements the processing manager for log entries.
package manager

import (
	"context"
	"errors"
	"io"
	"sync"

	"github.com/kubex-ecosystem/ethyr/tools/flow/control"
	"github.com/kubex-ecosystem/ethyr/tools/flow/events"
	"github.com/kubex-ecosystem/logz"
	kbxMod "github.com/kubex-ecosystem/xtui/internal/module/kbx"
)

// Manager handles the processing pipeline for log entries.
type Manager struct {
	mu sync.RWMutex

	formatter    Formatter
	writer       io.Writer
	hooks        events.Collection[logz.Entry]
	levelEnabled func(logz.Level) bool

	entry   logz.Entry
	options *kbxMod.InitArgs
	ctl     control.ManagerControl[control.StepFlag]
}

// Formatter defines the log formatting contract.
type Formatter interface {
	Format(e logz.Entry) ([]byte, error)
}

// IsTerminal checks if the manager is in a terminal state.
func (m *Manager) IsTerminal() bool {
	return m.ctl.IsTerminal()
}

// Process is the main pipeline for entries.
func (m *Manager) Process(ctx context.Context, entry logz.Entry) error {
	if entry == nil {
		return nil
	}

	if m.IsTerminal() {
		return control.ErrTerminal
	}

	// ---- Stage 1: validate/level gate --------------------------------------
	if m.levelEnabled != nil && !m.levelEnabled(logz.ParseLevel(entry.GetLevel())) {
		return nil
	}

	if err := m.stageValidate(entry); err != nil {
		m.advance(control.StepFailed)
		return err
	}
	m.advance(control.StepValidate)

	// ---- Stage 2: pre-hooks -------------------------------------------------
	if err := m.stagePreHooks(ctx, entry); err != nil {
		m.advance(control.StepFailed)
		return err
	}
	m.advance(control.StepPreHooks)

	// ---- Stage 3: format ----------------------------------------------------
	b, err := m.stageFormat(entry)
	if err != nil {
		m.advance(control.StepFailed)
		return err
	}
	m.advance(control.StepFormat)

	// ---- Stage 4: post-hooks ------------------------------------------------
	if err := m.stagePostHooks(ctx, entry); err != nil {
		m.advance(control.StepFailed)
		return err
	}
	m.advance(control.StepPostHooks)

	// ---- Stage 5: write -----------------------------------------------------
	if err := m.stageWrite(b); err != nil {
		m.advance(control.StepFailed)
		return err
	}
	m.advance(control.StepWrite)

	m.advance(control.StepDone)
	return nil
}

func (m *Manager) advance(flag control.StepFlag) {
	m.ctl.Stage.Store(flag)
}

func (m *Manager) stageValidate(entry logz.Entry) error {
	// Add custom validation if needed
	return nil
}

func (m *Manager) stagePreHooks(ctx context.Context, entry logz.Entry) error {
	return m.hooks.Fire(entry)
}

func (m *Manager) stageFormat(entry logz.Entry) ([]byte, error) {
	m.mu.RLock()
	f := m.formatter
	m.mu.RUnlock()

	if f == nil {
		return nil, errors.New("logz: no formatter configured in Manager")
	}

	b, err := f.Format(entry)
	if err != nil {
		return nil, err
	}

	if len(b) == 0 || b[len(b)-1] != '\n' {
		b = append(b, '\n')
	}

	return b, nil
}

func (m *Manager) stagePostHooks(ctx context.Context, entry logz.Entry) error {
	// Logic for post-hooks if different from pre-hooks
	return nil
}

func (m *Manager) stageWrite(b []byte) error {
	m.mu.RLock()
	out := m.writer
	m.mu.RUnlock()

	if out == nil {
		return errors.New("logz: no writer configured in Manager")
	}

	_, err := out.Write(b)
	return err
}

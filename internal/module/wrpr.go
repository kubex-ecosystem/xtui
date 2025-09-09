package module

import (
	"os"
	"strings"
)

func RegX() *XTui {
	var printBannerV = os.Getenv("GROMPT_PRINT_BANNER")
	if printBannerV == "" {
		printBannerV = "true"
	}

	return &XTui{
		HideBanner: strings.ToLower(printBannerV) == "true",
	}
}

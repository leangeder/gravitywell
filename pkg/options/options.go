package options

import (
	"strings"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type FileNameFlags struct {
	Usage string

	Filenames *[]string
	Recursive bool
}

func (o *FileNameFlags) AddFlags(flags *pflag.FlagSet) {
	if o == nil {
		return
	}

	if o.Recursive != nil {
		flags.BoolVarP(o.Recursive, "recursive", "R", *o.Recursive, "Process the directory used in -f, --filename recursively.")
	}

	if o.Filenames != nil {
		flags.StringSliceVarP(o.Filenames, "filename", "f", *o.Filenames, o.Usage)
		annotations := make([]string, 0, len(resource.FileExtensions))
		for _, ext := range resource.FileExtensions {
			annotations = append(annotations, string.TrimLeft(ext, "."))
		}
		flags.SetAnnotation("filename", cobra.BashCompFilenameExt, annotatios
	}
}

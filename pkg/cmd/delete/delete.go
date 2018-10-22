package delete

import "fmt"

type DeleteOptions struct {
	resource.FilenameOptions


}

func (o *DeleteOptions) RunDelete() error {
	return o.DeleteResult(o.Result)
}

func (o *DeleteOptions)

package file

import (
	"github.com/ignisVeneficus/library/api"
	"os"
)

func WriteAllBookToFile(filename string) error {
	data, err := api.GetAllBookAsJSON()
	if err != nil {
		return err
	}
	err = os.WriteFile("output.json", data, 0644)
	return err
}

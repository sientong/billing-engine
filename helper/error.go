package helper

import "log"

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func CheckErrorOrReturn(err error) error {
	if err != nil {
		log.Printf("Handled error: %v", err)
		return err
	}
	return nil
}

package valueutil

func CheckSetString(v *string, defValue string) {
	if len(*v) == 0 {
		*v = defValue
	}
}

func CheckSetInt64(v *int64, defValue int64) {
	if *v == 0 {
		*v = defValue
	}
}

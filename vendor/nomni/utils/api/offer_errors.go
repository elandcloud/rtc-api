package api

func RedundantItemsError(err error) Error {
	return newError(22001, "Redundant items", err.Error())
}

func InvalidTargetItemsError(err error) Error {
	return newError(22002, "Invalid target items", err.Error())
}

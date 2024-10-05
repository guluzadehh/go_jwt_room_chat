package api

func UnexpectedError() Response {
	return Err("an unexpecter error occured.")
}

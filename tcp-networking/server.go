package main

func server() error {
	endpoint := NewEndpoint()
	endpoint.AddHandlerFunc("STRING", handleStrings)
	endpoint.AddHandlerFunc("GOB", handleGob)
	return endpoint.Listen()
}

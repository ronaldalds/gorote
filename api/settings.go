package main

func Ready() error {
	return nil
}

func Config() error {
	Load()
	LoadSQL()
	if err := Ready(); err != nil {
		return err
	}
	return nil
}

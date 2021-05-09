package main

func authMethods() map[string]bool {
	path := "/aleg.laptops.LaptopService/"
	// SearchLaptop is accessible by everyone (even for
	// unregistered users).
	return map[string]bool{
		path + "CreateLaptop": true,
		path + "RateLaptop":   true,
		path + "UploadImage":  true,
	}
}

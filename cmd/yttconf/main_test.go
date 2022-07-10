package main

type AppStub struct{}

func (s *AppStub) Save() string {
	return "/home/john/.config/ytt/config.json"
}

func Example() {
	a = &AppStub{}
	main()

	// Output: Created /home/john/.config/ytt/config.json edit it with your favorite text editor
}

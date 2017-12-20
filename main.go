package main

func main() {
	var err error
	Client, err := NewClient()
	if err != nil {
		mlog.Errorln(err)
		return
	}
	State.Start()
	go Client.Query(domainSet)
	Client.Response()
	State.Show()
}

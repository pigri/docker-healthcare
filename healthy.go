package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/fredli74/lockfile"
)

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func configCheck() {

	if slackHookURL := os.Getenv("SLACK_HOOK_URL"); slackHookURL == "" {
		err := "SLACK_HOOK_URL environment variable is missing"
		panic(err)
	}

	if coreOsPrivateIpv4 := os.Getenv("COREOS_PRIVATE_IPV4"); coreOsPrivateIpv4 == "" {
		err := "COREOS_PRIVATE_IPV4 environment variable is missing"
		panic(err)
	}
}

func systemRun(cmd string) {
	out, err := exec.Command("/bin/sh", "-c", cmd).Output()
	checkError(err)
	fmt.Println(string(out))
}

func notification(status string) {
	webhookUrl := os.Getenv("SLACK_HOOK_URL")

	var color string

	switch status {
	case "timeout":
		color = "#ff0000"
	case "restarted":
		color = "#36a64f"
	}

	host := os.Getenv("COREOS_PRIVATE_IPV4")
	fallback := status + " on " + host

	attachment1 := slack.Attachment{Color: &color, Fallback: &fallback}
	attachment1.AddField(slack.Field{Title: "Host", Value: host}).AddField(slack.Field{Title: "Status", Value: status})
	payload := slack.Payload{
		Username:    "Docker healthy checker",
		Attachments: []slack.Attachment{attachment1},
	}

	err := slack.Send(webhookUrl, "", payload)
	if len(err) > 0 {
		fmt.Printf("error: %s\n", err)
	}
}

func main() {

	configCheck()

	if lock, err := lockfile.Lock("lockfile"); err != nil {
		checkError(err)
	} else {
		notification("timeout")
		systemRun("systemctl kill kubelet")
		systemRun("systemctl kill docker")
		systemRun("systemctl kill containerd")
		systemRun("systemctl start docker")
		systemRun("systemctl start kubelet")
		notification("restarted")
		defer lock.Unlock()
	}
}

package cmd

import (
	"testing"
	"fmt"
)

func TestRegister(t *testing.T) {
	fmt.Println("Testing Register...")
	got := Register("yanzexin", "123", "mail", "1234", "https://private-anon-f5f6a74e5f-agendaservice2.apiary-mock.com")
	expect := 0

	if got != expect {
		t.Errorf("got [%s] expected [%s]", got, expect)
	}
}

func TestLogin(t *testing.T) {
	fmt.Println("Testing Login...")
	got := Login("yanzexin", "123", "https://private-anon-f5f6a74e5f-agendaservice2.apiary-mock.com")
	expect := 0

	if got != expect {
		t.Errorf("got [%s] expected [%s]", got, expect)
	}
}

func TestLogout(t *testing.T) {
	fmt.Println("Testing Log out...")
	Login("yanzexin", "123", "https://private-anon-f5f6a74e5f-agendaservice2.apiary-mock.com")
	got := Logout("https://private-anon-f5f6a74e5f-agendaservice2.apiary-mock.com")
	expect := 0

	if got != expect {
		t.Errorf("got [%d] expected [%d]", got, expect)
	}
}

func TestShowInfo(t *testing.T) {
	fmt.Println("Testing show information...")
	Login("yanzexin", "123", "https://private-anon-f5f6a74e5f-agendaservice2.apiary-mock.com")
	got := ShowInfo("yanzexin", "https://private-anon-f5f6a74e5f-agendaservice2.apiary-mock.com")
	expect := 0

	if got != expect {
		t.Errorf("got [%d] expected [%d]", got, expect)
	}
}

func TestShowUsers(t *testing.T) {
	fmt.Println("Testing show users...")
	Login("yanzexin", "123", "https://private-anon-f5f6a74e5f-agendaservice2.apiary-mock.com")
	got := ShowUsers("https://private-anon-f5f6a74e5f-agendaservice2.apiary-mock.com")
	expect := 0

	if got != expect {
		t.Errorf("got [%d] expected [%d]", got, expect)
	}
}

func TestDeleteUser(t *testing.T) {
	fmt.Println("Testing delete user...")
	Login("yanzexin", "123", "https://private-anon-f5f6a74e5f-agendaservice2.apiary-mock.com")
	got := DeleteUser("yanzexin", "https://private-anon-f5f6a74e5f-agendaservice2.apiary-mock.com")
	expect := 0

	if got != expect {
		t.Errorf("got [%s] expected [%s]", got, expect)
	}
}

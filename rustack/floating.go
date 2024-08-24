package rustack

import (
	"fmt"
	"net/url"
)

type Floating struct {
	manager   *Manager
	ID        string `json:"id"`
	VDC       *Vdc   `json:"vdc,omitempty"`
	IpAddress string `json:"ip_address"`
}

func (m *Manager) GetFloating(id string) (fip *Floating, err error) {
	path, _ := url.JoinPath("v1/floating", id)
	err = m.Get(path, Defaults(), &fip)
	fip.manager = m
	return
}

func (v *Vdc) GetFloatingByAddress(address string) (fip *Floating, err error) {
	args := Arguments{
		"vdc":         v.ID,
		"filter_type": "external",
	}
	var items []*Floating
	err = v.manager.GetItems("v1/port", args, &items)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(items); i++ {
		if items[i].IpAddress == address {
			fip = items[i]
			fip.manager = v.manager
			return
		}
	}
	return nil, fmt.Errorf("ERROR. Address %s not found", address)
}

func (v *Vdc) CreateFloating() (*Floating, error) {
	args := &struct {
		VDC string `json:"vdc"`
	}{
		VDC: v.ID,
	}
	fip := new(Floating)
	err := v.manager.Request("POST", "v1/port", args, fip)
	if err != nil {
		return nil, err
	}
	fip.manager = v.manager
	return fip, nil
}

func (fip *Floating) Delete() error {
	err := fip.manager.Delete(fmt.Sprintf("v1/port/%s", fip.ID), Defaults(), nil)
	if err != nil {
		return err
	}
	return nil
}

func (fip Floating) WaitLock() (err error) {
	path := fmt.Sprintf("v1/port/%s", fip.ID)
	return loopWaitLock(fip.manager, path)
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

type APIServer struct{}

type ErrorResponse struct {
	Error string
}

type SuccessResponse struct {
	Message string
}

func (apiServer APIServer) Start() {
	var wg sync.WaitGroup
	log.Println("Starting API server...")
	http.HandleFunc("/", apiServer.Endpoint)
	wg.Add(1)
	go http.ListenAndServe(":8080", nil)
	log.Println("API Server started...")
	wg.Wait()
}

func (apiServer APIServer) Endpoint(w http.ResponseWriter, r *http.Request) {
	switch path := r.URL.Path[1:]; {
	case path == "deviceList" || path == "":
		apiServer.deviceList(w)
	case strings.HasPrefix(path, "turnOn/"):
		apiServer.turnOn(w, strings.Split(path, "turnOn/")[1])
	case strings.HasPrefix(path, "turnOff/"):
		apiServer.turnOff(w, strings.Split(path, "turnOff/")[1])
	case path == "health/live" || path == "health/ready":
		fmt.Fprintf(w, "ok")
	}
}

func getDevice(name string) (Device, ErrorResponse) {
	var device Device
	var err ErrorResponse
	for _, d := range config.Devices {
		if d.Name == name {
			return d, err
		}
	}
	err.Error = "Device not found."
	return device, err
}

func (apiServer APIServer) deviceList(w http.ResponseWriter) {
	var devices []Device
	for _, d := range config.Devices {
		d.Key = ""
		d.Status = ""
		devices = append(devices, d)
	}
	writeResponse(w, devices)
}

func (apiServer APIServer) turnOn(w http.ResponseWriter, deviceName string) {
	device, err := getDevice(deviceName)
	if err.Error != "" {
		writeResponse(w, err)
		return
	}
	turnOn(device.IP, device.Key)
	device.Status = "on"
	device.Key = ""
	writeResponse(w, device)
}

func (apiServer APIServer) turnOff(w http.ResponseWriter, deviceName string) {
	device, err := getDevice(deviceName)
	if err.Error != "" {
		writeResponse(w, err)
		return
	}
	turnOff(device.IP, device.Key)
	device.Status = "off"
	device.Key = ""
	writeResponse(w, device)
}

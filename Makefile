dir = /usr/bin
sysconfdir = /etc

build:
	@go build .

run: light

# made simple for later
install: light 90-backlight.rules
	install -CDt $(dir) light
	install -CDt $(sysconfdir)/udev/rules.d 90-backlight.rules

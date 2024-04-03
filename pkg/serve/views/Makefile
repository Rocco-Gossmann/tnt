TARGET := main.css 
SOURCE := main.tw.css

DEVFILES := $(wildcard *.tw.css) $(wildcard *.html)

remake:
	clear
	rm -rf $(TARGET)
	make $(TARGET)

$(TARGET): $(SOURCE)
	tailwindcss -i $< -o $@ -m --content './**/*.{html,tw.css}'

.PHONY: dev clean remake

dev:
	echo $(DEVFILES) | xargs -n 1 | entr make remake

clean:
	rm -rf $(TARGET)


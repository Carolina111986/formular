OUT := $(shell date +%s)
serve:
	ng serve --host 0.0.0.0

build:
	# mkdir ../../docs/$(OUT)
	ng build --output-path ../../docs
	git add ../../docs


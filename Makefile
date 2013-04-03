LIB?=	/usr/local/lib
BIN?=	/usr/local/bin
TCLSH?=	tclsh8.6

INSTALL_FILES=main.tcl lyrics.dat
PROGNAME=rickroll

all:

install:
	@echo Installing rickroll daemon
	install -o root -g wheel -m 0755 $(BIN)/tcllauncher $(BIN)/$(PROGNAME)
	install -o root -g wheel -m 0755 -d $(LIB)/$(PROGNAME)
	install -o root -g wheel -m 0644 $(INSTALL_FILES) $(LIB)/$(PROGNAME)

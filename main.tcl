#!/usr/local/bin/tclsh8.5

package require Tclx

proc chr_is_delimeter {chr} {
	if {$chr eq " " || $chr eq "-"} {
		return 1
	}
	return 0
}

proc waitfor {usec} {
	after $usec
}

proc main {} {
	set filename "lyrics.dat"
	set delay_word 200
	set delay_line 1000

	set filebuf [read_file [file join [file dirname [info script]] $filename]]

	foreach buf [split $filebuf "\n"] {
		set buf [string trim $buf]

		if {$buf eq ""} {
			puts "\r"
		} else {
			set lpos	0
			set buflen	[string length $buf]
			set lastchr [expr $buflen - 1]

			for {set rpos 0} {$rpos < $buflen} {incr rpos} {
				set chr [string range $buf $rpos $rpos]

				if {[chr_is_delimeter $chr] || $rpos == $lastchr} {
					if {[chr_is_delimeter $chr]} {
						set end [expr $rpos - 1]
					} else {
						set end $rpos
					}
					set word [string range $buf $lpos $end]

					if {[regexp {\{(\d+)\}} $word _ usec]} {
						waitfor $usec
						set last_output "delay"
					} else {
						puts -nonewline "$word"
						if {$chr eq " "} {
							puts -nonewline $chr
						}
						flush stdout
						set last_output "text"
						waitfor $delay_word
					}

					set lpos [expr $rpos + 1]
				}
			}
			puts "\r"

			if {$last_output ne "delay"} {
				waitfor $delay_line
			}
		}
	}

	return
}


if !$tcl_interactive main

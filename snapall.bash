#!/bin/bash

snapcraft clean
snapcraft
snapcraft push *.snap > push.txt
cat push.txt > push2.txt
grep 'Revision' push2.txt > push.txt
grep -o '[0-9]' push.txt > push2.txt
$(echo snapcraft release recipe-manager `tr -d "\n\r" < push2.txt` stable)
rm push.txt push2.txt

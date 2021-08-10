#!/bin/bash
for file in $(find /cddl-files -type f)
do
  echo "processing $file... "
  resultFile=${file//cddl/json}
  mkdir -p $(dirname "$resultFile")
  cddlc -tjson $file > $resultFile
  echo "ok!"
done

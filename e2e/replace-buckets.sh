#!/bin/bash
find . -type f -exec sed -i "s/BUCKET_1/$1/g" {} +
find . -type f -exec sed -i "s/BUCKET_2/$2/g" {} +

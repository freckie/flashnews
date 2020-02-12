import os
import time
import subprocess

subfolders = list()
for f in os.scandir('./'):
    if f.is_dir() and ('(keywords)' not in f.path):
        subfolders.append(f.path)

for it in subfolders:
    path = it + '/config.json'
    proc = subprocess.call('start ./main.exe "{}"'.format(path), shell=True)
    print('LOG: {} 폴더 실행.'.format(it))

print('일괄 실행 스크립트 종료.')
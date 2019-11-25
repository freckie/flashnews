import requests
from bs4 import BeautifulSoup

url = 'https://news.mt.co.kr/newsflash/newsflash.html?sec=all&listType=left'

req = requests.get(url)
bs = BeautifulSoup(req.content, 'lxml')

wrapper = bs.find('div', class_='group')
items = wrapper.find_all('li', class_='bundle')
print(items[0:10])
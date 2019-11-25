import requests
from bs4 import BeautifulSoup

url = 'https://www.asiae.co.kr/realtime/sokbo_left.htm'

req = requests.get(url)
bs = BeautifulSoup(req.content, 'lxml')

wrapper = bs.find('div', class_='ct txtform')
items = wrapper.find_all('li')
print(items[0:10])
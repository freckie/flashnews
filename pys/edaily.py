import requests
from bs4 import BeautifulSoup

url = 'https://www.edaily.co.kr/news/realtime/realtime_NewsList_1.asp'

req = requests.get(url)
bs = BeautifulSoup(req.content, 'lxml')

ul = bs.find('ul')
#items = [(it.find('span').get_text(), it.find('a').get_text()) for it in ul.find_all('li')]
for it in ul.find('li'):
    print(it)
#print(items)
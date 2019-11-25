import requests
from bs4 import BeautifulSoup

url = 'https://m.sedaily.com/News/NewsAll'

req = requests.get(url)
bs = BeautifulSoup(req.content, 'lxml')

ul = bs.find('ul', class_='news_list')
items = [it.get_text() for it in ul.find_all('a')]
print(items)

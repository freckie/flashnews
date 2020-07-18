import requests
from bs4 import BeautifulSoup

url = 'https://www.fnnews.com/load/category/002001000?page=0'

req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

items = bs.find_all('li')

for item in items[0:10]:
    a_tag = item.find('a')
    href = a_tag['href'].strip()
    title = item.find('strong', class_='tit_thumb').get_text().strip()
    date = item.find("em", class_='date').get_text().strip()

    req2 = requests.get(href)
    bs2 = BeautifulSoup(req2.text, 'lxml')
    
    wrapper2 = bs2.find('div', id='article_content')
    contents = ''.join([it.strip() for it in wrapper2.find_all(text=True, recursive=False)])

    print('=================')
    print(title)
    print(date)
    print(href)
    print(contents)
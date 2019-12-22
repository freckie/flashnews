import requests
from bs4 import BeautifulSoup

url = 'https://paxnetnews.com/categories'

req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

wrapper = bs.find('div', class_='list')
items = wrapper.find_all('div', class_='list-article')

for item in items[:1]:
    div = item.find('div', class_='content')
    a_tag = div.find('a', class_='dyn std')
    title = a_tag.get_text().strip()
    href = 'https://paxnetnews.com' + a_tag['href']
    date = item.find('div', class_='pubdate').get_text().strip()

    req2 = requests.get(href)
    bs2 = BeautifulSoup(req2.content, 'lxml')
    wrapper2 = bs2.find('article', class_='article')

    contents = ''.join([it.get_text().strip() for it in wrapper2.find_all('p')]).replace('\n', '')

    print(title, href, date, contents)
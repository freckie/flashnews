import requests
from bs4 import BeautifulSoup

url = 'https://www.yna.co.kr/news?site=navi_news'

req = requests.get(url)
bs = BeautifulSoup(req.content, 'lxml')

wrapper = bs.find('div', class_='headlines headline-list')
items = wrapper.find_all('li', class_='section02')

for item in items[:1]:
    dt = item.find('strong', class_='news-tl')

    a_tag = dt.find('a')
    title = a_tag.get_text().strip()
    href = 'https:' + a_tag['href']
    date = item.find('span', class_='p-time').get_text().strip()

    req2 = requests.get(href)
    bs2 = BeautifulSoup(req2.content, 'lxml')
    wrapper2 = bs2.find('div', class_='article')

    contents = ''.join([it.get_text().strip() for it in wrapper2.find_all('p')]).replace('\n', '')

    print(title, href, date, contents)
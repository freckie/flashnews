import requests
from bs4 import BeautifulSoup

url = 'http://inthenews.co.kr/latest-news/'
req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

wrapper = bs.find('div', class_='pt-cv-wrapper')
items = wrapper.find_all('div', class_='pt-cv-ifield')

for item in items[:4]:
    h4_tag = item.find('h4', class_='pt-cv-title')
    a_tag = h4_tag.find('a')
    href = a_tag['href']
    title = a_tag.get_text().strip()

    req2 = requests.get(href)
    bs2 = BeautifulSoup(req2.text, 'lxml')

    date = bs2.find('div', class_='post-date').get_text().strip()

    wrapper2 = bs2.find('div', class_='pf-content')
    contents = ''.join([it.get_text().strip() for it in wrapper2.find_all('p')])

    print(title, href, date, contents)
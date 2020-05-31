import requests
from bs4 import BeautifulSoup

url = 'http://m.kukinews.com/m/m_section.html?sec_no=66'
req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

wrapper = bs.find('ul', class_='lists')
items = wrapper.find_all('li')

# headline
headline = bs.find('div', class_='headline')
title = headline.find('a').get_text().strip()
url = headline.find('a')['href']

for item in items[:15]:
    a_tag = item.find('a')

    title = item.find('p', class_='tit').get_text().strip()
    url = 'http://m.kukinews.com' + a_tag['href']

    req2 = requests.get(url)
    bs2 = BeautifulSoup(req2.text, 'lxml')
    wrapper2 = bs2.find('div', id='news_body_area')

    date = bs2.find('div', class_='byline').get_text().split(' | ')[1].strip()
    contents = ''.join([it.get_text().strip() for it in wrapper2.find_all('p')])

    print(title)
    print(url)
    print(date)
    print(contents)
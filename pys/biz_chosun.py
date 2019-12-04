import requests
from bs4 import BeautifulSoup

url = 'https://biz.chosun.com/svc/bulletin/index.html'

req = requests.get(url)
bs = BeautifulSoup(req.content, 'lxml')

wrapper = bs.find('div', class_='art_list_wrap')
items = wrapper.find_all('li')

for item in items[:1]:
    a_tag = item.find('a')
    title = a_tag.get_text().strip()
    href = 'https://biz.chosun.com' + a_tag['href']
    date = item.find('span', class_='time').get_text().strip()
    title = title.replace(date, '')

    req2 = requests.get(href)
    bs2 = BeautifulSoup(req2.content, 'lxml')
    wrapper2 = bs2.find('div', id='news_body_id')
    title = bs2.find('h1', id='news_title_text_id').get_text().strip()

    contents = ''.join([it.get_text().strip() for it in wrapper2.find_all('div', class_='par')]).replace('\n', '')

    print(title, href, date, contents)
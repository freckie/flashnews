import requests
from bs4 import BeautifulSoup

url = 'http://news.mtn.co.kr/newscenter/news_section.mtn?scd=all'

req = requests.get(url)
bs = BeautifulSoup(req.content, 'lxml')

wrapper = bs.find('div', id='articleList')
items = wrapper.find_all('li')

for item in items[:1]:
    a_tag = item.find('h3').find('a')
    title = a_tag.get_text().strip()
    href = 'http://news.mtn.co.kr' + a_tag['href']
    date = item.find('span', class_='newsDate').get_text().strip()

    req2 = requests.get(href)
    bs2 = BeautifulSoup(req2.content, 'lxml')
    wrapper2 = bs2.find('div', id='newsContent')

    contents = ''.join([it.get_text().strip() for it in wrapper2.find_all('div')[:-4]]).replace('\n', '')

    print(title, href, date, contents)
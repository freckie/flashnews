import requests
from bs4 import BeautifulSoup

url = 'http://www.newspim.com/quicknews/left_list/?category_cd='

req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

wrapper = bs.find('ul', class_='newslist')
items = wrapper.find_all('li')

for item in items[:1]:
    a_tag = item.find('a')
    title = a_tag.find('span').get_text().strip()
    href = 'http://www.newspim.com' + a_tag['href']
    date = item.find('time').get_text()

    req2 = requests.get(href)
    bs2 = BeautifulSoup(req2.text, 'lxml')
    wrapper2 = bs2.find('div', id='news_contents')

    contents = ''.join([it.get_text().strip() for it in wrapper2.find_all('p')]).replace('\n', '')

    print(title, href, date, contents)
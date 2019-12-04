import requests
from bs4 import BeautifulSoup

url = 'http://www.fnnews.com/newsflash/'

req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

wrapper = bs.find('div', class_='art_flash_wrap')
items = wrapper.find_all('li')

for item in items[:1]:
    a_tag = item.find('a')
    title = a_tag.get_text().strip()
    href = 'http://www.fnnews.com' + a_tag['href']
    date = item.find('span', class_='category_date').get_text()

    req2 = requests.get(href)
    bs2 = BeautifulSoup(req2.text, 'lxml')
    wrapper2 = bs2.find('div', id='article_content')

    remove = wrapper2.find('table').get_text()
    contents = wrapper2.get_text().replace(remove, '').strip()

    print(title, href, date, contents)
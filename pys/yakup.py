import re
import requests
from bs4 import BeautifulSoup

url = 'https://www.yakup.com/news/index.html?cat=all'

req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

wrapper = bs.find('ul', class_='tp_5')
items = wrapper.find_all('div', class_='listBoxType_3')

for item in items[:1]:
    a_tag = item.find('a')
    title = a_tag.get_text().strip()
    href = 'https://www.yakup.com' + a_tag['href']
    span = item.find('p', class_='number tp_3').get_text()
    date = re.search('[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}', span).group()

    req2 = requests.get(href)
    bs2 = BeautifulSoup(req2.content, 'lxml')
    wrapper2 = bs2.find('div', class_='bodyarea')

    contents = ''.join([it.get_text().strip() for it in wrapper2.find_all('p')]).replace('\n', '')

    print(title, href, date, contents)
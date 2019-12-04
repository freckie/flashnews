import requests
from bs4 import BeautifulSoup

url = 'https://www.hankyung.com/all-news'

req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

wrapper = bs.find('ul', class_='article_list')
items = wrapper.find_all('li')

for item in items[:15]:
    a_tag = item.find('a')
    title = a_tag.get_text().strip()
    href = a_tag['href']
    date = item.find('p', class_='time').get_text()

    req2 = requests.get(href)
    bs2 = BeautifulSoup(req2.text, 'lxml')
    wrapper2 = bs2.find('div', id='articletxt')

    removes = wrapper2.find_all('div', class_='wrap_img')
    contents = wrapper2.get_text()
    for remove in removes:
        contents = contents.replace(remove.get_text(), '')

    contents = contents.replace('\n', '').replace('\t', '').replace('  ', '').strip()

    print(title, href, date, contents, end='\n\n')
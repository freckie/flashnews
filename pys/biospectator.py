import requests
from bs4 import BeautifulSoup

url = 'http://www.biospectator.com/section/section_list.php?MID=10000'

req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

wrapper = bs.find('div', class_='article_list')
items = wrapper.find_all('li')

for item in items[:1]:
    span = item.find('strong', class_='article_tit')
    a_tag = span.find('a')
    title = a_tag.get_text().strip()
    href = 'http://www.biospectator.com' + a_tag['href']
    date = item.find('span', class_='date').get_text().strip()

    req2 = requests.get(href)
    bs2 = BeautifulSoup(req2.content, 'lxml')
    wrapper2 = bs2.find('div', class_='article_view')

    p_tags = wrapper2.find_all('p')[1:]
    contents = ''
    for p in p_tags:
        if p.has_attr('class'):
            if p['class'][0] == 'photo_caption':
                continue
        contents = contents + p.get_text()
    print(title, href, date, contents)
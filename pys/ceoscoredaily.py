import requests
from bs4 import BeautifulSoup

url = 'http://www.ceoscoredaily.com/news/article_list_all.html'

req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

wrapper = bs.find('div', class_='article_list_002').find('ul')
items = wrapper.find_all('li')

result = list()
for item in items:
    temp = {
        'title': '',
        'link': '',
        'date': '',
        'content': ''
    }
    p_tag = item.find('p', class_='tit')
    a_tag = p_tag.find('a')

    temp['title'] = a_tag.get_text().strip()
    temp['link'] = 'http://www.ceoscoredaily.com/news/' + a_tag['href']
    temp['date'] = item.find('p', class_='date').get_text().strip()

    # contents
    req2 = requests.get(temp['link'])
    bs2 = BeautifulSoup(req2.text, 'lxml')
    wrapper2 = bs2.find('div', class_='article_body')
    #temp['contents'] = ' '.join([it.get_text().strip() for it in wrapper2.find_all('span', recursive=True)])
    temp['contents'] = wrapper2.get_text().strip().replace('\t', '').replace('\r', '').replace('\n', ' ')

    result.append(temp)

print(result)
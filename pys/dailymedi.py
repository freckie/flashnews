import requests
from bs4 import BeautifulSoup

url = 'http://www.dailymedi.com/ajax_section.php?ajaxNum=1&ajaxLayer=section_ajax_layer_1&thread=&numberpart=&file2=&file=normal_all_news.html&area=&pg=1&vals=%C0%DA%B5%BF%2C%C0%FC%C3%BC%2C%C3%D11%2F10%B0%B3%C3%E2%B7%C2%2C%C1%A6%B8%F142%C0%DA%C0%DA%B8%A7%2C%BA%BB%B9%AE250%C0%DA%C0%DA%B8%A7%2C%C5%F5%B8%ED%BB%F6%2C%B4%A9%B6%F40%B0%B3%2C%C0%FC%C3%BC%B4%BA%BD%BA%C3%E2%B7%C2%2C%C0%CC%B9%CC%C1%F6%B0%A1%B7%CE%C7%C8%BC%BF100%2Crows_photo_news07.html%2C%C0%DA%B5%BF%2C%C6%E4%C0%CC%C2%A1%2C&start_date=&end_date='

req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

wrapper = bs.find('table')
items = wrapper.find_all('tr')

for item in items[:1]:
    a_tag = item.find('a')
    href = 'http://www.dailymedi.com/' + a_tag['href']
    date = item.find('font').get_text().strip()

    req2 = requests.get(href)
    bs2 = BeautifulSoup(req2.content, 'lxml')
    wrapper2 = bs2.find('div', id='sub_center_contents2').find('table')
    title = wrapper2.find_all('tr', recursive=False)[3].find_all('tr')[1].get_text().strip()

    contents = wrapper2.find_all('tr', recursive=False)[6].find_all('table')[0].get_text().strip()

    print(title, href, date, contents)
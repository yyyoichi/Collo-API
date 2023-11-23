cd $HOME
wget https://github.com/shogo82148/mecab/releases/download/v0.996.8/mecab-0.996.8.tar.gz && \
  tar zxfv mecab-0.996.8.tar.gz && \
  cd mecab-0.996.8 && \
  ./configure && \
  make && \
  make check && \
  sudo make install && \
  sudo ldconfig && \
  # install unidic
  wget -P $HOME https://clrd.ninjal.ac.jp/unidic_archive/2302/unidic-cwj-202302.zip && \
  sudo mkdir -p /usr/local/lib/mecab/dic/unidic-cwj && \
  sudo unzip $HOME/unidic-cwj-202302.zip -d /usr/local/lib/mecab/dic/unidic-cwj && \
  # dic setting
  sudo sed -i 's#/usr/local/lib/mecab/dic/ipadic#/usr/local/lib/mecab/dic/unidic-cwj#g' /usr/local//etc/mecabrc

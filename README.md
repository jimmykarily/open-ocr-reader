*This project started as part of openSUSE HackWeek 21: https://hackweek.opensuse.org/21/projects/open-source-book-reader-for-visually-impaired*

**This is currently work in progress. It already produces some promising results but it's far from fulfilling its promise**

Planning board: https://github.com/jimmykarily/open-ocr-reader/projects/1

Reading a book is delight most of us can enjoy without thinking too much about it. It's not the same for blind or visually impaired people.
While most of the books are available in their traditional, paper form, not as many are available in digital form, let alone audio books.
There is a big variety of tools and applications that can help but they are either not complete, they offer low quality results, they
are expensive or not designed for bling people.

We have enough technology available around us to solve this problem in an optimal way. This project doesn't try to invent the wheel but
rather collect and combine existing solutions to this problem. The goal is to allow blind people to read a book as easily as possible.
Ideally, it should be possible to put a book in front of a camera, click a button and seconds later the book should be read to the user
in a voice as natural as possible

## Goals of this project

- Make reading books for blind and visually impaired people as easy as possible
- Make is as cheap as possible (free is possible)
- Use free and open source tools as much as possible
- Offer more options when free/open source tools are not sufficient.

## Non-goals of this project

- Create new OCR libraries
- Create new Text to speech libraries

## Architecture

The process to get from text on paper, to audio is described in the following image:

![architecture](assets/architecture.svg)

## Alternatives to this project

- https://www.readforme.io/ (source code?)

## Useful links:

- Page frame detection: https://users.iit.demokritos.gr/~bgat/CBDAR_BORDERS.pdf

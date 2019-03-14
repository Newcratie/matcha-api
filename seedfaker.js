[
    '{{repeat(5, 7)}}',
    {
        username: '{{firstName()}}{{integer(0, 1000)}}',
        email: '{{email()}}',
        lastname: '{{surname()}}',
        firstname: '{{firstName()}}',
        password: '1234567890',
        created_at: '{{date(new Date(2014, 0, 1), new Date(), "YYYY-MM-ddThh:mm:ss Z")}}',
        random_token: '{{objectId()}}',
        img1: 'https://randomuser.me/api/portraits/men/{{integer(1,40)}}.jpg',
        img2: 'https://randomuser.me/api/portraits/men/{{integer(1,40)}}.jpg',
        img3: 'https://randomuser.me/api/portraits/men/{{integer(1,40)}}.jpg',
        img4: 'https://randomuser.me/api/portraits/men/{{integer(1,40)}}.jpg',
        img5: 'https://randomuser.me/api/portraits/men/{{integer(1,40)}}.jpg',
        biography: '{{lorem(1, "paragraphs")}}',
        birthday: '{{date(new Date(2014, 0, 1), new Date(), "YYYY-MM-ddThh:mm:ss Z")}}',
        genre: 'male',
        interest: '{{random("hetero","homo","bi")}}',
        city: '{{city()}}',
        zip: '{{integer(1000, 97000)}}',
        country: '{{state()}}',
        latitude: '{{floating(-90.000001, 90)}}',
        longitude: '{{floating(-180.000001, 180)}}',
        online: '{{bool()}}',
        latitude: '{{floating(-90.000001, 90)}}',
        longitude: '{{floating(-180.000001, 180)}}',
        geo_allowed: '{{bool()}}',
        rating: '{{floating(0, 10)}}',
        admin: '{{bool()}}'
    }

]
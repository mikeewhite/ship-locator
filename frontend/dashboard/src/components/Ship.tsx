import React from 'react';
import { Descriptions, Space } from 'antd';
import type { DescriptionsProps } from 'antd';

const Ship = (props: { ship: any; }) => {
    const { ship } = props;

    const borderedItems: DescriptionsProps['items'] = [
        {
            key: '1',
            label: 'MMSI',
            children: ship.mmsi,
        },
        {
            key: '2',
            label: 'Name',
            children: ship.name,
        },
        {
            key: '3',
            label: 'Latitude',
            children: ship.latitude,
        },
        {
            key: '4',
            label: 'Longitude',
            children: ship.longitude,
        },
    ]

    return (
        <Space direction="vertical">
            <br/>
            <Descriptions
                bordered
                column={1}
                size="small"
                items={borderedItems}
            />
            <iframe width="600" height="450" frameBorder="0" style={{ border: 0}}
                    src={"https://www.google.com/maps/embed/v1/view?key=" + process.env.REACT_APP_GOOGLE_MAPS_API_KEY + "&center=" + ship.latitude + "," + ship.longitude + "&zoom=15"}
                    allowFullScreen>
            </iframe>
        </Space>
    );
};

export default Ship;
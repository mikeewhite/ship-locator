import React, { useState } from 'react';
import { Empty, Space } from 'antd';
import { useLazyQuery, gql } from '@apollo/client';
import Ship from './Ship';
import ShipAutoComplete from "./ShipAutoComplete";

const SHIP_SEARCH_QUERY = gql`
    query ($filter: Int!) {
        ship(mmsi: $filter) {
            mmsi
            name
            latitude
            longitude
            lastUpdated
        }
    }
`;

const ShipSearch = () => {
    const [selectedData, setSelectedData] = useState<String | null>(null);
    const [searchFilter, setSearchFilter] = useState('');
    const [executeSearch, { data }] = useLazyQuery(
        SHIP_SEARCH_QUERY
    );

    const handleAutoCompleteSelect = (mmsi: String | null) => {
        executeSearch({
            variables: { filter: mmsi }
        }).then(r => {
                if (r.data) {
                    setSelectedData(mmsi);
                }
            }
        );
    };

    return (
        <Space direction="vertical">
            <ShipAutoComplete onSelect={handleAutoCompleteSelect}/>
                {selectedData ? (
                            <Ship ship={data.ship} />
                    ) : (
                        <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} />
                    )
                }
        </Space>
    );
};

export default ShipSearch;
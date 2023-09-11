import React, { useState } from 'react';
import { Empty, Input } from 'antd';
import { useLazyQuery, gql } from '@apollo/client';
import Ship from './Ship';

const { Search } = Input;

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
    const [searchFilter, setSearchFilter] = useState('');
    const [executeSearch, { data }] = useLazyQuery(
        SHIP_SEARCH_QUERY
    );

    return (
        <>
            <div>
                <Search placeholder="Enter MMSI..."
                        style={{ width: 200 }}
                        onChange={(e) => setSearchFilter(e.target.value)}
                        onSearch={() => executeSearch({
                            variables: { filter: searchFilter }
                        })}
                />
            </div>
            {data ? (
                    <Ship ship={data.ship} />
                ) : (
                    <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} />
                )
            }
        </>
    );
};

export default ShipSearch;
import React from 'react';
import { Layout, Space, version, theme } from 'antd';
import { AimOutlined, RocketOutlined } from '@ant-design/icons';
import ShipSearch from "./ShipSearch";

const { Header, Footer, Content } = Layout;

function App() {
    const {
        token: { colorBgContainer },
    } = theme.useToken();

    return (
      <Space direction="vertical" style={{ width: '100%' }}>
          <Layout className="layout" style={{ minHeight: '100vh' }}>
              <Header style={{ display: 'flex', alignItems: 'center' }}>
                  <div className="demo-logo"><h1 style={{ color: 'white' }}><AimOutlined/> Ship Locator</h1></div>
              </Header>
              <Layout style={{ padding: '24px 0', background: colorBgContainer }}>
                  <Content style={{ textAlign: 'center', padding: '0 24px' }}>
                      <ShipSearch/>
                  </Content>
              </Layout>
              <Footer style={{ textAlign: 'center' }}><RocketOutlined/> Built with Ant Design v{version}</Footer>
          </Layout>
      </Space>
  );
}

export default App;

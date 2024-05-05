/**
 * 主应用组件
 * 该组件负责渲染整个应用的界面结构，使用了Ant Design组件库和react-router-dom库来进行布局和路由管理。
 */
import React from 'react';
import { Layout, Menu, Typography } from 'antd';
import { CarryOutOutlined, UserOutlined, FileOutlined, FundOutlined } from '@ant-design/icons';
import Job from './pages/Job'
import EditJob from './pages/EditJob'
import Account from './pages/Account'
import EditAccount from './pages/EditAccount'
import EditTemplate from './pages/EditTemplate'
import Template from './pages/Template'
import Log from './pages/Log'
import 'antd/dist/reset.css';
import './index.css';
import { HashRouter, Route, Link, Routes } from 'react-router-dom';
import { globalConfig } from "./config";
import { createRoot } from 'react-dom/client';

// 引入 AntDesign组件
const { Header, Content, Footer, Sider } = Layout;
const { Title } = Typography;

const container = document.getElementById('root');
const root = createRoot(container); // 创建根容器

/**
 * 应用主类
 */
class App extends React.Component {
    state = {
        collapsed: false,
    };

     /**
     * 渲染函数，负责渲染整个应用界面
     * @returns 返回应用的JSX结构
     */
    onCollapse = collapsed => {
        this.setState({ collapsed });
    };

    render() {
        // 设置主页标题
        document.title = globalConfig.rootTitle;
        return (
            <HashRouter>
                <Layout style={{ minHeight: '100vh' }}>
                    {/*左侧导航条*/}
                    <Sider collapsible collapsed={this.state.collapsed} onCollapse={this.onCollapse}>
                        <Title level={4} style={{color:'white', margin:10, textAlign:'center'}}>监控小子 (Surveillance-guy)</Title>
                        <Menu theme="dark" defaultSelectedKeys={[window.location.hash]} mode="inline">
                            <Menu.Item key="#/job" link='/job'>
                                <CarryOutOutlined/>
                                <Link to='/job'>定时任务</Link>
                            </Menu.Item>
                            <Menu.Item key="#/account" link='/account'>
                                <UserOutlined/>
                                <Link to='/account'>通知账户</Link>
                            </Menu.Item>
                            <Menu.Item key="#/template">
                                <FileOutlined/>
                                <Link to='/template'>模板配置</Link>
                            </Menu.Item>
                            <Menu.Item key="#/log">
                                <FundOutlined/>
                                <Link to='/log'>实时日志</Link>
                            </Menu.Item>
                        </Menu>
                    </Sider>
                    <Layout>
                        <Header style={{ background: '#fff', padding: 0, height: 51, position: 'relative' }} />
                        {/*    <Button style={{ position: 'absolute', top: '20%', marginLeft: '16px' }} onClick={this.goBack}>*/}
                        {/*        <span>返回</span>*/}
                        {/*    </Button>*/}
                        {/*</Header>*/}
                        {/*右侧显示内容*/}
                        <Content style={{ margin: '16px 16px' }}>
                            {/*面包屑*/}
                            {/*<Breadcrumb >*/}
                            {/*</Breadcrumb>*/}
                            {/*正文*/}
                            <div style={{ padding: 24, minHeight: 360, background: '#fff' }}>
                                {/*路由配置，根据路由显示不同的页面内容*/}
                                <Routes>
                                    <Route path='/' breadcrumbName="首页" element={<Job/>} />
                                    <Route path='/job' breadcrumbName="定时任务" element={<Job/>} />
                                    <Route path='/editjob' breadcrumbName="创建任务" element={<EditJob/>} />
                                    <Route path='/account' breadcrumbName="通知账户" element={<Account/>} />
                                    <Route path='/editaccount' breadcrumbName="添加账户" element={<EditAccount/>} />
                                    <Route path='/template' breadcrumbName="模板配置" element={<Template/>} />
                                    <Route path='/edittemplate' breadcrumbName="创建模板" element={<EditTemplate/>} />
                                    <Route path='/log' breadcrumbName="实时日志" element={<Log/>} />
                                </Routes>
                            </div>

                        </Content>
                        {/*页脚 底部信息*/}
                        <Footer style={{ textAlign: 'center' }}>©2024 Edited by 711lxsky</Footer>
                    </Layout>
                </Layout>
            </HashRouter>
        );
    }
}

root.render(<App />);

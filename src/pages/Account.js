/**
 * Account组件用于展示和管理邮件账户信息
 * 包括账户列表的展示，账户操作（编辑、测试、删除），以及账户状态的实时更新。
 */
import React from 'react';
import { Table, Button, message, Radio, Popconfirm, Badge, Tag, Tooltip } from 'antd';
import { Link, useNavigate, } from 'react-router-dom';
import axios from 'axios';
import { globalConfig } from '../config'
import PropTypes from 'prop-types';
import { InfoCircleOutlined, LoadingOutlined } from '@ant-design/icons';

function withNavigation(Component) {
    return function WrapperComponent(props) {
        const navigate = useNavigate();
        return <Component {...props} navigate={navigate} />;
    };
}

class Account extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            accounts: [], // 邮件账户列表
        };
    }

    // 表格列的定义
    columns = [
        {
            title: 'Email', // 列标题
            dataIndex: 'email', // 数据索引
            key: 'email', // 唯一标识
            width: '30%', // 列宽度
            // 自定义渲染内容，用于处理过长的文本
            render: (text) => (
                <div style={{ wordWrap: 'break-word', wordBreak: 'break-all' }}>
                    {text}
                </div>
            ),
        },
        {
            title: '密码 / 授权码',
            dataIndex: 'password',
            key: 'password',
            // 显示为星号，保护敏感信息
            render: () => ('********'),
        },
        {
            // 自定义复杂列，包含是否可用的状态展示和提示信息
            title: <Tooltip placement="right" title="可以点击【测试】按钮检测 Email 账号是否可用"> 是否可用 <InfoCircleOutlined theme="twoTone" /></Tooltip>,
            key: 'status',
            dataIndex: 'status',
            render: (text) => {
                // 根据状态显示不同的徽标
                switch (text) {
                    case 1:
                        return <Badge status="success" text={<Tag color="green">可用</Tag>} />;
                    case 2:
                        return <Badge status="error" text={<Tag color="red">不可用</Tag>} />;
                    case 3:
                        return <Badge status="processing" text={<Tag color="blue">测试中</Tag>} />;
                    default:
                        return <Badge status="default" text={<Tag color="gray">未测试</Tag>} />;
                }
            },
        },
        {
            title: '操作',
            key: 'action',
            // 自定义操作列，包含编辑、测试和删除操作
            render: (_text, record) => {
                return (
                    <span>
                        <Radio.Group>
                            <Link to={{ pathname: '/editaccount', state: { account: record } }}>
                                <Radio.Button>编辑</Radio.Button>
                            </Link>
                            <Radio.Button onClick={() => { this.testEmail(record) }} >
                                {record.status === 3 && <LoadingOutlined style={{ marginRight: 5 }} />}
                                测试
                            </Radio.Button>
                            <Popconfirm
                                title="是否确定删除？"
                                onConfirm={() => { this.handleDelete(record) }}
                                okText="是"
                                cancelText="否"
                            >
                                <Radio.Button>删除</Radio.Button>
                            </Popconfirm>
                        </Radio.Group>
                    </span>
                )
            },
        },
    ];

    /**
     * 测试当前Email账户的连通性
     * @param {Object} record - 要测试的账户记录
     */
    testEmail = (record) => {
        // 找到对应记录的索引
        const i = this.state.accounts.findIndex(item => record.ID === item.ID);

        // 将账户设置为测试中状态，并使用回调确保正确引用前一个状态来更新state
        this.setState(prevState => ({
            accounts: prevState.accounts.map((account, index) =>
                index === i ? { ...account, status: 3 } : account
            )
        }), () => this.forceUpdate()); // 这里的forceUpdate可能不是必要的，因为setState通常会导致重新渲染

        const data = {
            id: record.ID,
            email: record.email,
            host: record.host,
            port: parseInt(record.port),
        };

        // 发送测试请求
        axios.post(globalConfig.rootPath + '/api/v1/testemail', data)
            .then(() => {
                // 成功时，同样使用回调更新状态
                this.setState(prevState => {
                    const updatedAccounts = [...prevState.accounts]; // 复制数组以避免直接修改状态
                    updatedAccounts[i].status = 1;
                    return { accounts: updatedAccounts };
                }, () => {
                    message.success("该 Email 账户身份验证通过，可以发送邮件");
                });
            })
            .catch(e => {
                // 失败时，也利用setState的回调功能更新状态
                this.setState(prevState => {
                    const failedAccounts = [...prevState.accounts];
                    failedAccounts[i].status = 2;
                    return { accounts: failedAccounts };
                }, () => {
                    console.log(e);
                    if (e.response && e.response.data) {
                        if (e.response.data.message) {
                            message.error(`[消息] ${e.response.data.message} [原因] ${e.response.data.reason}`);
                        } else {
                            // 如果没有 message 属性，可以提供默认错误信息
                            message.error('发生未知错误');
                        }
                    } else {
                        // 如果没有 response 对象，可能是网络错误或其他问题
                        message.error(e.message || '网络错误');
                    }

                });
            });
    };

    /**
     * 组件将要挂载时调用，用于同步账户数据
     */
    componentDidMount() {
        this.syncAccounts();
    }

    /**
     * 处理账户删除操作
     * @param {Object} record - 要删除的账户记录
     */
    handleDelete = (record) => {
        // 判断是否满足至少有一个通知账户
        if (this.state.accounts.length <= 1) {
            message.warning("至少应有一个通知账户");
            return;
        }
        // 发送删除请求
        axios.delete(globalConfig.rootPath + '/api/v1/account', { data: JSON.stringify(record) })
            .then(res => {
                console.log(res);
                if (res.status === 200) {
                    message.success('删除成功');
                    // 更新state
                    this.setState((prevState) => {
                        const afterAccounts = prevState.accounts.filter((v) => v.ID !== record.ID);
                        return { accounts: afterAccounts };
                    });
                    this.setState({ 'accounts': afterAccounts });
                    // 检查删除后是否至少有一个账户
                    this.validateAtLeastOneEmail(afterAccounts);
                }
            })
            .catch(e => {
                console.log(e);
                if (e?.e.response?.e.response.data?.e.response.data.message)
                    message.error("[message] " + e.response.data.message + " [reason] " + e.response.data.reason);
                else
                    message.error(e.message);
            });
    };

    /**
     * 同步账户数据到state
     */
    syncAccounts = () => {
        axios.get(globalConfig.rootPath + '/api/v1/account')
            .then(res => {
                console.log(res);
                let accounts = res.data.data;
                this.setState({ 'accounts': accounts });
                // 检查同步后的账户数量
                this.validateAtLeastOneEmail(accounts);
            })
            .catch(e => {
                console.log(e);
                if (e?.response?.data?.message)
                    message.error("[message] " + e.response.data.message + " [reason] " + e.response.data.reason);
                else
                    message.error(e.message);
            });
    };

    static propTypes = {
        navigate: PropTypes.func.isRequired, // 声明navigate prop的类型
    };

    /**
     * 检查是否至少有一个邮件账户，如果没有则跳转到添加账户页面
     * @param {Array} accounts - 账户列表
     */
    validateAtLeastOneEmail = (accounts) => {
        if (accounts.length === 0) {
            message.warning("请至少添加一个通知账户");
            this.props.navigate('/editaccount');
        }
    };

    render() {
        // 渲染表格和添加账户按钮
        return (
            <div>
                <Table
                    rowKey="ID"
                    columns={this.columns}
                    dataSource={this.state.accounts}
                    pagination={false}
                    bordered
                />

                <Link to='/editaccount'>
                    <Button
                        style={{ width: '100%', margin: '16px 0', height: 40 }}
                        icon="plus"
                    >
                        添加账户
                    </Button>
                </Link>
            </div>
        )
    }
}

export default withNavigation(Account);
import React, { useState, useEffect } from 'react';
import { Form, Input, Button, message, Typography } from 'antd';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import axios from 'axios';
import { globalConfig } from '../config';
import { MailOutlined, LockOutlined, DatabaseOutlined, ApiOutlined } from '@ant-design/icons';

const { Title } = Typography;

function EditAccount() {
    const navigate = useNavigate();
    const location = useLocation();
    const [form] = Form.useForm();
    const [isEdit, setIsEdit] = useState(false);
    const [account, setAccount] = useState(null);
    const [testEmailStatus, setTestEmailStatus] = useState("play-circle");

    useEffect(() => {
        if (location.state && location.state.account) {
            setIsEdit(true);
            setAccount(location.state.account);
        }
    }, [location.state]);

    const handleSubmit = async () => {
        try {
            const values = await form.validateFields();
            values.port = parseInt(values.port);
            let method = isEdit ? axios.put : axios.post;
            if (isEdit) values.ID = account.ID;

            const res = await method(globalConfig.rootPath + '/api/v1/account', values);
            console.log(res);
            if (res.status === 200) {
                message.success(isEdit ? '通知账户更新成功' : '通知账户创建成功');
                navigate(-1);
            }
        } catch (error) {
            console.error('Failed to submit form:', error);
        }
    };

    const testEmail = async () => {
        setTestEmailStatus("loading");
        try {
            const email = form.getFieldValue("email");
            const password = form.getFieldValue("password");
            const host = form.getFieldValue("host");
            const port = parseInt(form.getFieldValue("port"));

            let data = {
                id: account && account.ID ? account.ID : 0,
                email,
                password,
                host,
                port,
            };

            await axios.post(globalConfig.rootPath + '/api/v1/testemail', data);
            setTestEmailStatus("check-circle");
            message.success("该 Email 账户身份验证通过，可以发送邮件");
        } catch (error) {
            console.error('Failed to test email:', error);
            setTestEmailStatus("close-circle");
        }
    };

    return (
        <Form form={form} onFinish={handleSubmit}>
            <Title level={4} type="secondary"><Link to="#" onClick={() => navigate(-1)} style={{ marginRight: '5px' }}>通知账户配置</Link></Title>
            <Form.Item
                label="Email"
                name="email"
                rules={[
                    {
                        type: 'email',
                        message: '不是有效的 Email',
                    },
                    {
                        required: true,
                        message: '请输入邮件账号 (Email)',
                    },
                ]}
                initialValue={isEdit? account.email : ''}
            >
                <Input
                    prefix={<MailOutlined style={{ color: 'rgba(0,0,0,.25)' }} />}
                    placeholder="请输入邮件账号..."
                />
            </Form.Item>
            <Form.Item
                label="密码 / 授权码"
                name="password"
                rules={[
                    { required: true, message: '请输入密码 / 授权码' },
                ]}
                extra={
                    <Button
                        type="primary"
                        size="small"
                        icon={testEmailStatus}
                        onClick={testEmail}
                        danger={testEmailStatus === "close-circle"}
                        ghost
                    >
                        测试连通性
                    </Button>
                }
            >
                <Input.Password
                    prefix={<LockOutlined style={{ color: 'rgba(0,0,0,.25)' }} />}
                    placeholder="请输入密码 / 授权码..."
                />
            </Form.Item>
            {/* SMTP 服务器地址表单项 */}
            <Form.Item
                label="SMTP 服务器地址（可选）"
                name="host"
                rules={[
                    { optional: true, pattern: /\S+/u, message: '请输入 SMTP 服务器地址' },
                ]}
                initialValue={isEdit ? account.host : ''}
            >
                <Input
                    prefix={<DatabaseOutlined style={{ color: 'rgba(0,0,0,.25)' }} />}
                    placeholder="请输入 SMTP 服务器地址..."
                />
            </Form.Item>

            {/* SMTP 服务器端口表单项 */}
            <Form.Item
                label="SMTP 服务器端口（可选）"
                name="port"
                rules={[
                    { optional: true, type: 'number', min: 1, max: 65535, message: '请输入有效的 SMTP 服务器端口' },
                ]}
                initialValue={isEdit && account.port !== 0 ? account.port.toString() : ''}
            >
                <Input
                    prefix={<ApiOutlined style={{ color: 'rgba(0,0,0,.25)' }} />}
                    placeholder="请输入 SMTP 服务器端口..."
                    type="number" // 添加 input 类型为 number
                />
            </Form.Item>

            <Form.Item>
                <Button type="primary" htmlType="submit">
                    {isEdit ? '提交' : '添加'}
                </Button>
                <Link to='/account'>
                    <Button style={{ marginLeft: '10px' }}>
                        返回
                    </Button>
                </Link>
            </Form.Item>
        </Form>
    );
}

export default EditAccount;

import React, { useState, useEffect } from 'react';
import {
    Form,
    Input,
    Select,
    Button,
    message,
    Typography,
    Divider,
    Drawer,
    Dropdown,
    Menu
} from 'antd';
import { InputCron , LeftOutlined} from 'antcloud-react-crons';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import axios from 'axios';
import ReactMarkdown from 'react-markdown';
import { globalConfig } from '../config';

const { Option } = Select;
const { TextArea } = Input;

function EditJob() {
    const navigate = useNavigate();
    const location = useLocation();
    const [form] = Form.useForm();
    const [isVisibleDrawer, setIsVisibleDrawer] = useState(false);
    const [emails, setEmails] = useState([]);
    const [templates, setTemplates] = useState([]);
    const [job, setJob] = useState(location.state?.job || {});
    const [testPatternStatus, setTestPatternStatus] = useState("play-circle");

    useEffect(() => {
        getAllAccounts();
        getAllTemplates();
        if (location.state && location.state.job) {
            setJob(location.state.job);
        }
    }, [location.state]);

    const getAllAccounts = async () => {
        try {
            const response = await axios.get(`${globalConfig.rootPath}/api/v1/account`);
            const accounts = response.data.data;
            setEmails(accounts.map(account => account.email));
            if (accounts.length === 0) {
                message.warning("请至少添加一个通知账户");
                navigate('/editaccount');
            }
        } catch (error) {
            message.error("获取账户信息失败");
        }
    };

    const getAllTemplates = async () => {
        try {
            const response = await axios.get(`${globalConfig.rootPath}/api/v1/template`);
            setTemplates(response.data.data);
        } catch (error) {
            message.error("获取模板信息失败");
        }
    };

    const testPattern = async () => {
        try {
            setTestPatternStatus("loading");
            const values = await form.validateFields(['url', 'pattern']);
            const params = {
                id: job?.ID || 0,
                url: values.url,
                type: "re",
                pattern: values.pattern,
            };
            const response = await axios.get(`${globalConfig.rootPath}/api/v1/testpattern`, { params });
            setTestPatternStatus("check-circle");
            message.success("匹配结果 => " + response.data.data);
        } catch (error) {
            setTestPatternStatus("close-circle");
            message.error("测试模式失败");
        }
    };

    const handleSubmit = async (values) => {
        try {
            const method = job.ID ? axios.put : axios.post;
            const response = await method(`${globalConfig.rootPath}/api/v1/job`, { ...values, ID: job.ID });
            message.success(job.ID ? '定时任务更新成功' : '定时任务创建成功');
            navigate(-1);
        } catch (error) {
            message.error("提交任务失败");
        }
    };

    const handleTemplateClick = e => {
        const selectedTemplate = templates.find(template => template.name === e.key);
        if (selectedTemplate) {
            form.setFieldsValue({
                cron: selectedTemplate.cron,
                pattern: selectedTemplate.pattern,
                content: selectedTemplate.content
            });
        }
    };

    // 定义 onShowDrawer 函数，用于显示抽屉
    const onShowDrawer = () => {
        setIsVisibleDrawer(true);
    };

    // 定义 onCloseDrawer 函数，用于隐藏抽屉
    const onCloseDrawer = () => {
        setIsVisibleDrawer(false);
    };

    return (
        <div id="job-form">
            <Form form={form} onFinish={handleSubmit} layout="vertical">
                <Typography.Title level={4}>
                    <Link to="#" onClick={() => navigate(-1)} style={{ marginRight: '5px' }}>
                        <LeftOutlined />
                    </Link>
                    定时任务配置
                </Typography.Title>
                <Divider />
                <Form.Item name="name" label="任务名称" initialValue={job?.name} rules={[{ required: true, message: '请输入任务名称' }]}>
                    <Input />
                </Form.Item>
                <Form.Item name="cron" label="定时配置" initialValue={job?.cron} rules={[{ required: true, message: '请输入定时配置' }]}>
                    <InputCron lang='zh_CN' />
                </Form.Item>
                <Form.Item name="url" label="目标页面 URL" initialValue={job?.url} rules={[{ required: true, message: '请输入目标页面 URL' }, { type: 'url', message: '请输入有效的URL' }]}>
                    <Input />
                </Form.Item>
                <Form.Item name="pattern" label="抓取规则" initialValue={job?.pattern} rules={[{ required: true, message: '请输入抓取规则' }]}>
                    <TextArea rows={4} />
                </Form.Item>
                <Form.Item name="email" label="通知账户" initialValue={job?.email} rules={[{ required: true, message: '请选择通知账户 Email' }]}>
                    <Select>
                        {emails.map(email => (
                            <Option key={email} value={email}>{email}</Option>
                        ))}
                    </Select>
                </Form.Item>
                <Form.Item name="content" label="邮件内容" initialValue={job?.content} rules={[{ required: true, message: '请输入邮件内容' }]}>
                    <TextArea rows={4} />
                </Form.Item>
                {/* 其他 Form.Item 组件 */}
                <Form.Item>
                    <Button type="primary" htmlType="submit">
                        {job ? '更新' : '创建'}
                    </Button>
                    <Link to="/job">
                        <Button style={{ marginLeft: '10px' }}>返回</Button>
                    </Link>
                </Form.Item>
            </Form>

            <Drawer title="正则表达式使用手册" placement="right" width={800} onClose={onCloseDrawer} visible={isVisibleDrawer}>
                <ReactMarkdown>{/* Markdown content here */}</ReactMarkdown>
            </Drawer>
        </div>
    );
}

export default EditJob;